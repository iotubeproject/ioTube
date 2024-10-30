records=(
  "'0x124B468860ac994ef28D11643d16e88dd4F26085' '0x99B2B0eFb56E62E36960c20cD5ca8eC6ABD5557A' 1 '0x7AD34903Bae921D44e5aF398383cc48e374847d4' ''"
  "'0x540a92Dd951407Ee6c94b997a43ecF30Ea6D04Cd' '0x99B2B0eFb56E62E36960c20cD5ca8eC6ABD5557A' 4592 '0x9517D9B64b2A17CF3b0f48Eff381ae556f9e2852' ''"
)
for record in "${records[@]}"; do
  echo $record
  IFS=" " read -r cashier token tidx recipient payload <<< "$record"
  docker exec witness-db mysql -uroot -pkdfjjrU64fjK58H -e "INSERT IGNORE INTO witness.iotex_to_matic_transfers_shadow (cashier, token, tidx, recipient, payload) VALUES (${cashier}, $token, $tidx, $recipient, $payload)"
  docker exec witness-db mysql -uroot -pkdfjjrU64fjK58H -e "UPDATE witness.iotex_to_matic_transfers SET status='ready' WHERE status='confirmed' AND cashier=$cashier AND token=$token AND tidx=$tidx"
done
