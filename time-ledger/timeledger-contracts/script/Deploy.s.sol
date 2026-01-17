// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/TimeLedgerToken.sol";

contract Deploy is Script {
    function run() external {
        // 从 .env 读取私钥
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");

        // 开始广播交易
        vm.startBroadcast(deployerPrivateKey);

        // 部署合约
        TimeLedgerToken token = new TimeLedgerToken("TimeLedger Token", "TLT");

        vm.stopBroadcast();

        // 打印合约地址
        console2.log("TimeLedgerToken deployed at:", address(token));
    }
}
