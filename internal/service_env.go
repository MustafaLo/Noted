package internal
import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnv()(map[string]string, error){
	data, err := godotenv.Read()
	if err != nil{
		return nil, fmt.Errorf("error loading env file: %w", err)
	}

	return data, nil
}

func UpdateEnv(new_env map[string]string)(error){
	if err := godotenv.Write(new_env, ".env"); err != nil {
		return fmt.Errorf("error writing to .env file: %w", err)
	}
	return nil
}