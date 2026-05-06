package repository

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DbMigrator gorm.Migrator

func ConnectToDatabase(host string, port string, user string, dbName string, password string) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", host, port, user, dbName, password)
	// fmt.Printf("host=%v port=%v user=%v dbname=%v password=%v", host, port, user, dbName, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	DB = db
	DbMigrator = db.Migrator()
}

var embedMigrations embed.FS

func RunMigrations() {
	// Migrate the schema
	goose.SetBaseFS(embedMigrations)
	sqlDb, err := DB.DB()
	if err != nil {
		panic(err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlDb, "../../migrations"); err != nil {
		panic(err)
	}
}

func CheckTableValue() error {
	if !DbMigrator.HasTable("users") {
		log.Fatal("Table 'users' does not exist. Please run the migrations.")
	}

	ctx := context.Background()

	rtvPok, err := gorm.G[models.PokemonCard](DB).Find(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(rtvPok) == 0 {
		fmt.Println("No Pokémon cards found in the database. Please run the seeders.")
		return RunSeeder("pokemon")
	}

	rtvNaruto, err := gorm.G[models.NarutoCard](DB).Find(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(rtvNaruto) == 0 {
		fmt.Println("No Naruto cards found in the database. Please run the seeders.")
		return RunSeeder("naruto")
	}

	rtvDigimon, err := gorm.G[models.DigimonCard](DB).Find(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(rtvDigimon) == 0 {
		fmt.Println("No Digimon cards found in the database. Please run the seeders.")
		return RunSeeder("digimon")
	}

	return nil
}

func RunSeeder(flag string) error {
	log.Println("Running seeder with flag:", flag)
	cmd := exec.Command("go", "run", "../../cmd/seeder/seeder.go", "-mode ", flag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()

}
