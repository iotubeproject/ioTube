#!/bin/bash
# Test cases for witness signing scripts

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Test private key (DO NOT USE IN PRODUCTION)
TEST_PRIVATE_KEY="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
TEST_SECRET_FILE="/tmp/test-witness-secret.yaml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

pass_count=0
fail_count=0

pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((pass_count++))
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    ((fail_count++))
}

section() {
    echo ""
    echo -e "${YELLOW}=== $1 ===${NC}"
}

# Setup
setup() {
    echo "privateKey: \"$TEST_PRIVATE_KEY\"" > "$TEST_SECRET_FILE"
}

cleanup() {
    rm -f "$TEST_SECRET_FILE"
}

trap cleanup EXIT
setup

# ============================================
# Test 1: sign-witness binary exists
# ============================================
section "Test 1: Binary Existence"

if [[ -x "./sign-witness" ]]; then
    pass "sign-witness binary exists and is executable"
else
    fail "sign-witness binary not found or not executable"
fi

if [[ -x "./sign-tx" ]]; then
    pass "sign-tx binary exists and is executable"
else
    fail "sign-tx binary not found or not executable"
fi

if [[ -x "./submit-witness" ]]; then
    pass "submit-witness binary exists and is executable"
else
    fail "submit-witness binary not found or not executable"
fi

# ============================================
# Test 2: sign-witness help
# ============================================
section "Test 2: Help Output"

if ./sign-witness -h 2>&1 | grep -q "Generates a witness signature"; then
    pass "sign-witness -h shows help"
else
    fail "sign-witness -h help output missing"
fi

if ./sign-to-eth.sh 2>&1 | grep -q "Usage"; then
    pass "sign-to-eth.sh shows usage without args"
else
    fail "sign-to-eth.sh usage output missing"
fi

if ./sign-to-solana.sh 2>&1 | grep -q "Usage"; then
    pass "sign-to-solana.sh shows usage without args"
else
    fail "sign-to-solana.sh usage output missing"
fi

if ./sign-tx.sh 2>&1 | grep -q "Usage"; then
    pass "sign-tx.sh shows usage without args"
else
    fail "sign-tx.sh usage output missing"
fi

if ./submit-witness.sh 2>&1 | grep -q "Usage"; then
    pass "submit-witness.sh shows usage without args"
else
    fail "submit-witness.sh usage output missing"
fi

# ============================================
# Test 3: sign-witness manual mode (Ethereum)
# ============================================
section "Test 3: Ethereum Signature Generation"

OUTPUT=$(./sign-witness -secret "$TEST_SECRET_FILE" \
    -cashier-address 0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3 \
    -validator-address 0xEf503971Aec1BC3cF3D896742Fa82975dCcB3162 \
    -cotoken 0x8bb487Db47c9F957CBd995AedDFbaa8895522F3D \
    -index 123 \
    -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
    -recipient 0x1234567890123456789012345678901234567890 \
    -amount 1000000000000000000 2>&1)

if echo "$OUTPUT" | grep -q "Transfer ID:"; then
    pass "sign-witness generates Transfer ID"
else
    fail "sign-witness missing Transfer ID"
fi

if echo "$OUTPUT" | grep -q "Signature:"; then
    pass "sign-witness generates Signature"
else
    fail "sign-witness missing Signature"
fi

# Verify signature is 65 bytes (130 hex chars + 0x prefix)
SIG=$(echo "$OUTPUT" | grep "Signature:" | awk '{print $2}')
SIG_LEN=${#SIG}
if [[ $SIG_LEN -eq 132 ]]; then  # 0x + 130 hex chars
    pass "Ethereum signature is 65 bytes"
else
    fail "Ethereum signature wrong length: $SIG_LEN (expected 132)"
fi

# ============================================
# Test 4: sign-witness manual mode (Solana)
# ============================================
section "Test 4: Solana Signature Generation (Ed25519)"

OUTPUT=$(./sign-witness -to-solana -secret "$TEST_SECRET_FILE" \
    -cashier-address 0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3 \
    -validator-address 0xEf503971Aec1BC3cF3D896742Fa82975dCcB3162 \
    -cotoken EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v \
    -index 456 \
    -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
    -recipient 9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM \
    -amount 1000000 2>&1)

if echo "$OUTPUT" | grep -q "To Solana: true"; then
    pass "sign-witness detects Solana mode"
else
    fail "sign-witness Solana mode not detected"
fi

# Verify Ed25519 signature is 64 bytes (128 hex chars + 0x prefix)
SIG=$(echo "$OUTPUT" | grep "Signature:" | awk '{print $2}')
SIG_LEN=${#SIG}
if [[ $SIG_LEN -eq 130 ]]; then  # 0x + 128 hex chars
    pass "Solana signature is 64 bytes (Ed25519)"
else
    fail "Solana signature wrong length: $SIG_LEN (expected 130)"
fi

# ============================================
# Test 5: sign-to-eth.sh with config
# ============================================
section "Test 5: sign-to-eth.sh with Config"

if [[ -f "configs/witness-config-iotex-testnet.yaml" ]]; then
    OUTPUT=$(./sign-to-eth.sh \
        -config configs/witness-config-iotex-testnet.yaml \
        -secret "$TEST_SECRET_FILE" \
        -cashier iotex-testnet-to-bsc-testnet \
        -token 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
        -index 789 \
        -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
        -recipient 0x1234567890123456789012345678901234567890 \
        -amount 500000000 2>&1)

    if echo "$OUTPUT" | grep -q "Transfer ID:"; then
        pass "sign-to-eth.sh works with config"
    else
        fail "sign-to-eth.sh failed with config"
    fi

    # Verify it extracted the correct co-token
    if echo "$OUTPUT" | grep -q "0x8bb487db47c9f957cbd995aeddfbaa8895522f3d"; then
        pass "sign-to-eth.sh extracts correct co-token from config"
    else
        fail "sign-to-eth.sh co-token extraction failed"
    fi
else
    fail "sign-to-eth.sh config file not found"
fi

# ============================================
# Test 6: sign-to-solana.sh with config
# ============================================
section "Test 6: sign-to-solana.sh with Config"

if [[ -f "configs/witness-config-iotex-solana.yaml" ]]; then
    OUTPUT=$(./sign-to-solana.sh \
        -config configs/witness-config-iotex-solana.yaml \
        -secret "$TEST_SECRET_FILE" \
        -cashier iotex-to-solana \
        -token io158elyywekvljpp4gqzzzsdkk0ufehxn6g3aa0u \
        -index 999 \
        -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
        -recipient 9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM \
        -amount 1000000 2>&1)

    if echo "$OUTPUT" | grep -q "Transfer ID:"; then
        pass "sign-to-solana.sh works with config"
    else
        fail "sign-to-solana.sh failed with config"
    fi
else
    fail "sign-to-solana.sh config file not found"
fi

# ============================================
# Test 7: Error handling
# ============================================
section "Test 7: Error Handling"

# Missing required parameter
if ./sign-witness -secret "$TEST_SECRET_FILE" -cotoken 0x8bb487Db47c9F957CBd995AedDFbaa8895522F3D 2>&1 | grep -q "required"; then
    pass "sign-witness reports missing required params"
else
    fail "sign-witness should report missing params"
fi

# Invalid cashier ID
if ./sign-to-eth.sh \
    -config configs/witness-config-iotex-testnet.yaml \
    -secret "$TEST_SECRET_FILE" \
    -cashier invalid-cashier-id \
    -token 0x8bb487Db47c9F957CBd995AedDFbaa8895522F3D \
    -index 1 -sender 0x... -recipient 0x... -amount 1 2>&1 | grep -q "Error"; then
    pass "sign-to-eth.sh reports invalid cashier ID"
else
    fail "sign-to-eth.sh should report invalid cashier ID"
fi

# ============================================
# Test 8: Deterministic signature
# ============================================
section "Test 8: Signature Determinism"

# Same input should produce same output
OUTPUT1=$(./sign-witness -secret "$TEST_SECRET_FILE" \
    -cashier-address 0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3 \
    -validator-address 0xEf503971Aec1BC3cF3D896742Fa82975dCcB3162 \
    -cotoken 0x8bb487Db47c9F957CBd995AedDFbaa8895522F3D \
    -index 123 \
    -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
    -recipient 0x1234567890123456789012345678901234567890 \
    -amount 1000000000000000000 2>&1 | grep "Transfer ID:" | awk '{print $3}')

OUTPUT2=$(./sign-witness -secret "$TEST_SECRET_FILE" \
    -cashier-address 0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3 \
    -validator-address 0xEf503971Aec1BC3cF3D896742Fa82975dCcB3162 \
    -cotoken 0x8bb487Db47c9F957CBd995AedDFbaa8895522F3D \
    -index 123 \
    -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
    -recipient 0x1234567890123456789012345678901234567890 \
    -amount 1000000000000000000 2>&1 | grep "Transfer ID:" | awk '{print $3}')

if [[ "$OUTPUT1" == "$OUTPUT2" ]]; then
    pass "Transfer ID is deterministic"
else
    fail "Transfer ID should be deterministic"
fi

# ============================================
# Test 9: submit-witness dry run
# ============================================
section "Test 9: Submit Witness Dry Run"

# Create a fake signature (65 bytes)
FAKE_SIG="0x$(printf '%0130d' 0)"

OUTPUT=$(./submit-witness \
    -config configs/relayer-config-iotex-testnet.yaml \
    -secret "$TEST_SECRET_FILE" \
    -cashier 0x5E0Eba3f0c9e047BbcD88441865F97643ab97Fd3 \
    -token 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
    -index 123 \
    -sender 0xFd57f47E48eC422599Fa44c4F370D7a474B38bBb \
    -recipient 0x1234567890123456789012345678901234567890 \
    -amount 1000000000000000000 \
    -signatures "$FAKE_SIG" \
    -dry-run 2>&1)

if echo "$OUTPUT" | grep -q "DRY RUN"; then
    pass "submit-witness dry-run works"
else
    fail "submit-witness dry-run failed"
fi

# ============================================
# Summary
# ============================================
echo ""
echo "============================================"
echo -e "${GREEN}Passed: $pass_count${NC}"
echo -e "${RED}Failed: $fail_count${NC}"
echo "============================================"

if [[ $fail_count -eq 0 ]]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
