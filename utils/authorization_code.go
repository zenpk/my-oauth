package utils

type AuthorizationInfo struct {
	ClientId      string
	Uuid          string
	CodeChallenge string
}

var AuthorizationCodeMap map[string]AuthorizationInfo

func GenAuthorizationCode(info AuthorizationInfo) (string, error) {
	code, err := RandString(Conf.AuthorizationCodeLength)
	if err != nil {
		return "", err
	}
	AuthorizationCodeMap[code] = info
	return code, nil
}
