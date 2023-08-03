package utils

const (
	HttpAddress    = "0.0.0.0:4430"
	InvitationCode = "code"

	AccessTokenAge = 24 // hour
	Issuer         = "github.com/zenpk/my-oauth"
	JwtPrivateKey  = "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
	JwtPublicKey   = `{
  "kty": "RSA",
  "n": "tvQSq9qpYrZGyK_fFFexPUHH-vkSRtqi4rk61hSbGdKuBqRDHruiAPhVvXRGP0ilqsBHtHKY1tKOf06mYWjbmKPvJtv7CitFRG5c86E5G30PqTd9QMXhsROCUSLdHwt-lMcSNsJSJBJ0RDNrdxgpBkajDY8AClkRUVFhMTjvQnq8RTVRIJ96WBz8k5T-D-lj5kz_2-ROhp9K5xDsSb_emAOxBNf2-EzqiWsDn2GXecApKYMW3pe1_cZUialjPwk8MGfY401Y1UC924PoQJOz4wd99jcNs4cHvaBIqxsUHylApguQjufLTI5bt8d9pY40xC540_jOlt4q77wD3zG1OUt9aMCA7of8K4mlI7PZweMOntwh7Shp0LaQQd_vzZ7iKPfIKA0W1hUQ4loSGnOC54wZGi2AFyjV-Zop_l_JZr53f7JIp3LqJGn3QCKChrOjCDGri9mHH_LfKBb2RQ-Hv5_ljStdS-A4cng3MX-ANxRN2Q7cSk1mH_C-VwT7wb-kPk2z2ilGrIZ_rxqssp6ioidVwXQMDNEvbyMLntYKOtYCFk7afPows-XfgyTZwDWx5pTC3pYQDrzl-gL_gSMoVNOOUgQOcj4md2GJ3IdTBZPHufPXIPcE2G12ssYmGhy3Tu1U1y3Zw-i6-K63182DJM969aUq9OlCp0T2Sn4nPR0",
  "e": "AQAB",
  "alg": "RS256",
  "kid": "123",
  "use": "sig"
}`
)
