package libUtils

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/types"
	"github.com/spf13/viper"
)

func GetConnConfigFromViper(v *viper.Viper) (types.ConnConfig, error) {
	connConfig := types.ConnConfig{}
	parseErr := viper.UnmarshalKey(constants.DbConfigConnSubkey, &connConfig)
	if parseErr != nil {
		return connConfig, parseErr
	}

	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)
	connConfig.Database = currBranch

	return connConfig, nil
}
