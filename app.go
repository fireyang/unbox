package unbox

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "app"}

func init() {
	/*
		var cmd = &cobra.Command{
			Use:   "hello",
			Short: "hello",
			Args: func(cmd *cobra.Command, args []string) error {
				if len(args) < 1 {
					return errors.New("requires at least one arg")
				}
				// return fmt.Errorf("invalid color specified: %s", args[0])
				fmt.Println("test")
				return nil
			},
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hello, World!")
			},
		}
		cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
		rootCmd.AddCommand(cmd)
	*/

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(unpackCmd)
	unpackCmd.Flags().StringVarP(&AtlasType, "type", "t", "json", "atlas type")
	unpackCmd.Flags().StringVarP(&OutPath, "out", "o", "./out", "out path")
}

var AtlasType string
var OutPath string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Unbox",
	Long:  `Print the version number of Unbox`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Unbox 0.1")
	},
}

var unpackCmd = &cobra.Command{
	Use:   "atlas",
	Short: "upack atlas",
	Long:  `unpack atlas`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}

		if _, err := os.Stat(OutPath); os.IsNotExist(err) {
			err = os.MkdirAll(OutPath, 0777)
			if err != nil {
				return err
			}
		}

		if AtlasType != "json" {
			return errors.New("other type come soon")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var inputPath = args[0]
		fi, err := os.Stat(inputPath)
		if err != nil {
			fmt.Printf("err:%v \n", err)
			return
		}
		ext := "." + AtlasType
		switch mode := fi.Mode(); {
		case mode.IsDir():
			fmt.Println("directory")
			err := UnpackDir(inputPath, OutPath, ext)
			if err != nil {
				fmt.Printf("err:%v \n", err)
			}
		case mode.IsRegular():
			// fmt.Printf("file:%v,%v\n", path.Ext(inputPath), AtlasType)
			if path.Ext(inputPath) != ext {
				fmt.Printf("ext err \n")
				return
			}
			err := UnpackFile(inputPath, OutPath)
			if err != nil {
				fmt.Printf("err:%v \n", err)
			}
			fmt.Println("unpack success")
		}
	},
}
