generate-ecdsa:
	openssl ecparam -name prime256v1 -genkey -noout -out private.pem