package cmd

import (
	"fmt"
	"github.com/mehgokalp/insider-project/internal/domain/mysql"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"math/rand"
)

func PopulateCmd(db *gorm.DB) *cobra.Command {
	cmdName := "populate"

	return &cobra.Command{
		Use:   cmdName,
		Short: "Populate database with dummy data",
		RunE: func(cmd *cobra.Command, _ []string) error {
			for i := 0; i < 20; i++ {
				msg := mysql.Message{
					To:      fmt.Sprintf("+90555%06d", rand.Intn(999999)),
					Content: fmt.Sprintf("Test message %d", i+1),
					Sent:    false,
				}
				db.Create(&msg)
			}

			return nil
		},
	}
}
