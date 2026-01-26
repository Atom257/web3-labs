package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
)

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

// Register registers all HTTP routes
func (s *Server) Register(r *gin.Engine) {
	r.GET("/head", s.GetHead)

	r.GET("/user/points", s.GetUserPoints)
	r.GET("/user/point_logs", s.GetUserPointLogs)

	r.GET("/rate/current", s.GetCurrentRate)
}

//
// =======================
// Handlers
// =======================
//

// GET /head?chain_id=&contract=
func (s *Server) GetHead(c *gin.Context) {
	chainID, err := strconv.ParseInt(c.Query("chain_id"), 10, 64)
	if err != nil || chainID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chain_id"})
		return
	}
	contract := c.Query("contract")
	if contract == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing contract"})
		return
	}

	type HeadResp struct {
		ChainID         int64  `json:"chain_id"`
		ContractAddress string `json:"contract_address"`
		Finalized       struct {
			BlockNumber int64  `json:"block_number"`
			BlockHash   string `json:"block_hash"`
		} `json:"finalized"`
	}

	var bh models.BlockHeader
	err = s.db.
		Where("chain_id=? AND contract_address=?", chainID, contract).
		Order("block_number DESC").
		Limit(1).
		First(&bh).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "head not found"})
		return
	}

	var resp HeadResp
	resp.ChainID = chainID
	resp.ContractAddress = contract
	resp.Finalized.BlockNumber = bh.BlockNumber
	resp.Finalized.BlockHash = bh.BlockHash

	c.JSON(http.StatusOK, resp)
}

// GET /user/points?chain_id=&contract=&account=
func (s *Server) GetUserPoints(c *gin.Context) {
	chainID, err := strconv.ParseInt(c.Query("chain_id"), 10, 64)
	if err != nil || chainID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chain_id"})
		return
	}

	contract := c.Query("contract")
	account := c.Query("account")
	if contract == "" || account == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing contract or account"})
		return
	}

	// UserPoint 目前是单一大表，不需要分表逻辑
	var up models.UserPoint
	err = s.db.
		Where(
			"chain_id=? AND contract_address=? AND account=?",
			chainID, contract, account,
		).
		First(&up).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user point not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chain_id":         up.ChainID,
		"contract_address": up.ContractAddress,
		"account":          up.Account,
		"total_points":     up.TotalPoints,
		"last_calc_time":   up.LastCalcTime,
	})
}

// GET /user/point_logs?chain_id=&contract=&account=&from=&to=&limit=&offset=
func (s *Server) GetUserPointLogs(c *gin.Context) {
	chainID, err := strconv.ParseInt(c.Query("chain_id"), 10, 64)
	if err != nil || chainID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chain_id"})
		return
	}

	contract := c.Query("contract")
	account := c.Query("account")
	if contract == "" || account == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing contract or account"})
		return
	}

	limit := 100
	offset := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	// 1. 先查询 SysContract 获取 ID，从而确定表名
	var sysC models.SysContract
	if err := s.db.Where("chain_id = ? AND address = ?", chainID, contract).First(&sysC).Error; err != nil {
		// 如果合约都找不到，肯定没有日志
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not configured"})
		return
	}

	// 获取动态表名 (例如 user_point_log_1)
	logTableName := sysC.GetLogTableName()

	// 2. 使用 Table(logTableName) 查询
	q := s.db.
		Table(logTableName).
		Where(
			"chain_id=? AND contract_address=? AND account=?",
			chainID, contract, account,
		)

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			q = q.Where("from_time >= ?", t)
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			q = q.Where("to_time <= ?", t)
		}
	}

	var logs []models.UserPointLog
	if err := q.
		Order("from_time ASC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {

		// 如果表不存在（可能是刚添加合约还没算过分），返回空列表或报错均可
		// 这里选择返回 500 暴露问题，或者你可以 check error 类型忽略 "table not found"
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GET /rate/current?chain_id=&contract=&at=
func (s *Server) GetCurrentRate(c *gin.Context) {
	chainID, err := strconv.ParseInt(c.Query("chain_id"), 10, 64)
	if err != nil || chainID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chain_id"})
		return
	}

	contract := c.Query("contract")
	if contract == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing contract"})
		return
	}

	at := time.Now().UTC()
	if v := c.Query("at"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			at = t
		}
	}

	var pr models.PointRate
	err = s.db.
		Where(
			"chain_id=? AND contract_address=? AND effective_time <= ?",
			chainID, contract, at,
		).
		Order("effective_time DESC").
		Limit(1).
		First(&pr).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rate not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chain_id":         pr.ChainID,
		"contract_address": pr.ContractAddress,
		"rate_numerator":   pr.RateNumerator,
		"rate_denominator": pr.RateDenominator,
		"effective_time":   pr.EffectiveTime,
	})
}
