package cmd

import "log"

func exitWithError(action string, err error) {
	if err != nil {
		log.Fatalf("failed to %s due to error %v", action, err)
	}
}
