package configuration

type Configuration struct {
	EnvironmentVariables *EnvironmentVariables
}

func UpConfiguration() (*Configuration, error) {

	environmentVariables, err := loadEnvVariables()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		EnvironmentVariables: environmentVariables,
	}, nil
}
