package main
import("context";"os";"os/signal";"syscall";"github.com/A1b3rt0M3rcad0/woobe-cli/internal/cli")
func main(){ctx,stop:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM);defer stop();os.Exit(cli.New(os.Stdin,os.Stdout,os.Stderr).Execute(ctx,os.Args[1:]))}
