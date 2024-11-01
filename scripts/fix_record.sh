records=(
  "'0x124B468860ac994ef28D11643d16e88dd4F26085' '0x99B2B0eFb56E62E36960c20cD5ca8eC6ABD5557A' 2 '0x5EE71bdEed6af2ea3e826C4e0250026e1029f031' ''"
  "'0x124B468860ac994ef28D11643d16e88dd4F26085' '0x8e66c0d6B70C0B23d39f4B21A1eAC52BBA8eD89a' 1 '0x8C5984dFa1f1C1709f881c79Cf70c6e24054B986' ''"
  "'0x124B468860ac994ef28D11643d16e88dd4F26085' '0x3CDb7c48E70B854ED2Fa392E21687501D84B3AFc' 1 '0x2319fB270317Ed97132f6e35Ce831544d93A3920' ''"
)
docker stop iotex-witness
docker exec witness-db mysql -uroot -pkdfjjrU64fjK58H -e "UPDATE witness.iotex_to_matic_transfers SET status='ready' WHERE status <> 'settled' AND cashier='0x540a92Dd951407Ee6c94b997a43ecF30Ea6D04Cd'"
for record in "${records[@]}"; do
  echo $record
  IFS=" " read -r cashier token tidx recipient payload <<< "$record"
  docker exec witness-db mysql -uroot -pkdfjjrU64fjK58H -e "INSERT IGNORE INTO witness.iotex_to_matic_transfers_shadow (cashier, token, tidx, recipient, payload) VALUES (${cashier}, $token, $tidx, $recipient, $payload)"
  docker exec witness-db mysql -uroot -pkdfjjrU64fjK58H -e "UPDATE witness.iotex_to_matic_transfers SET status='ready' WHERE status='confirmed' AND cashier=$cashier AND token=$token AND tidx=$tidx"
done
sleep 120s
docker start iotex-witness
