Setup Vault

1. Run vault-compose.yml: docker-compose -f vault-compose.yml up


2. Check Vault Status: docker exec vault vault status
   - Vault must be sealed and uninitialized 


3. Initialize Vault: docker exec -it vault vault operator init
   - Save the unseal keys and root token in the .env file
   - Add `VAULT_ADDR="http://127.0.0.1:8200"` to the .env file


4. Unseal Vault: docker exec -it vault vault operator unseal <unseal_key>
   - Repeat this step for 3 unseal keys


5. Login to Vault: docker exec -it vault vault login <root_token>


6. Enable KV Secrets Engine: docker exec -it vault vault secrets enable -path=secret kv


7. Write Secret: docker exec -it vault vault kv put secret/jwt_token token="my-super-secret-token"


8. Read Secret: docker exec -it vault vault kv get secret/jwt_token


9. Finally, you can access the secret in your application by sending a request to the Vault API. 
