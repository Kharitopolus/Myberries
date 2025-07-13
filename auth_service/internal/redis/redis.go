package redis

import "github.com/redis/go-redis/v9"

type RClient struct {
	client                   *redis.Client
	refreshToeknExpTimeHOurs int
}

func New(url string, refreshTokenExpTimeHours int) RClient {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	return RClient{client: redis.NewClient(opt), refreshToeknExpTimeHOurs: refreshTokenExpTimeHours}
}
