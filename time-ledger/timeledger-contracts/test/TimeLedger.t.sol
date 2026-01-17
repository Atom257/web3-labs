// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/TimeLedgerToken.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
contract TimeLedgerTest is Test {
    TimeLedgerToken token;

    address owner = address(this);
    address alice = address(0xA11CE);
    address bob = address(0xB0B);

    event Transfer(address indexed from, address indexed to, uint256 value);

    function setUp() public {
        token = new TimeLedgerToken("TimeLedger Token", "TLT");
    }
    function testMintBurnTransfer() public {
        // 初始 Alice 0
        assertEq(token.balanceOf(alice), 0);

        // mint 100 给 Alice
        token.mint(alice, 100 ether);
        assertEq(token.balanceOf(alice), 100 ether);

        // 再 mint 100
        token.mint(alice, 100 ether);
        assertEq(token.balanceOf(alice), 200 ether);

        // Alice 转 50 给 Bob
        vm.prank(alice);
        token.transfer(bob, 50 ether);

        assertEq(token.balanceOf(alice), 150 ether);
        assertEq(token.balanceOf(bob), 50 ether);

        // burn Alice 的 20
        token.burn(alice, 20 ether);
        assertEq(token.balanceOf(alice), 130 ether);

        // burn Bob 的 20
        token.burn(bob, 20 ether);
        assertEq(token.balanceOf(bob), 30 ether);
    }

    function testTransferEvents() public {
        // 期待 mint 100 给 Alice
        vm.expectEmit(true, true, false, true);
        emit Transfer(address(0), alice, 100 ether);
        token.mint(alice, 100 ether);

        // 期待再 mint 100
        vm.expectEmit(true, true, false, true);
        emit Transfer(address(0), alice, 100 ether);
        token.mint(alice, 100 ether);

        // Alice -> Bob 转 50
        vm.prank(alice);
        vm.expectEmit(true, true, false, true);
        emit Transfer(alice, bob, 50 ether);
        token.transfer(bob, 50 ether);

        // burn 20
        vm.expectEmit(true, true, false, true);
        emit Transfer(alice, address(0), 20 ether);
        token.burn(alice, 20 ether);
    }
}
