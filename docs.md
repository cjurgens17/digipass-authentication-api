uint32 allows for 2^32 || 4,294,967,295 seconds = 49,710 days more than enough for a grant limit of lets say one year, 2^16 is to low

86400 seconds in a day

Steps to creating JWT:
1. Create the Header
    1. Choose ALG (Usually HMAC SHA256 || RSA)
    2. Serialize to JSON
    3. Encode Base64Url
2. Create the Payload (Add private Claims if not null)
    1. Generate All Public Claim Values
    2. Get Private Claim values
    3. Serialize To JSON
    4. Encode Base64Url
3. Create the Signature
    1. Take Encoded Header
    2. Take Encoded Payload
    3. Take Secret
    4. Use Alg in Header and sign
4. Store the jti with the jwt for grant identification
5. Return JWT