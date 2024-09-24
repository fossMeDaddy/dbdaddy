package libUtils

import "github.com/fossmedaddy/dbdaddy/constants"

func GetDbConfigOriginKey(originName string) string {
	return constants.DbConfigOriginsKey + "." + originName
}
