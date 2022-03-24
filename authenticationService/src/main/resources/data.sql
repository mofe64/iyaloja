INSERT INTO oauth_client_details
(client_id, client_secret, web_server_redirect_uri, scope, access_token_validity, refresh_token_validity, resource_ids,
 authorized_grant_types, additional_information)
VALUES ('iyaloja', '{bcrypt}$2a$12$EerwEsmkUnDPLzbKtZ7mL.iMaRTzdI56r5ZOA66dcKYaf2Sm9xZaW', 'http://localhost:8080/',
        'internal,public', '3600', '10000', 'inventory,payment', 'authorization_code,password,refresh_token,implicit', '{}');
