package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.genframe.xyz/config"
)

type Bootstrap struct {
	Config     *config.Config
	CurrentDir string

	genframeProc    *os.Process
	autoupdaterProc *os.Process
	//
	notify chan error
}

func ProvideBootstrap(c *config.Config) *Bootstrap {
	path, _ := os.Getwd()

	return &Bootstrap{
		Config:     c,
		CurrentDir: path,
		notify:     make(chan error, 1),
	}
}

func (b *Bootstrap) StartGenframe() error {
	genframeProc, err := os.StartProcess(fmt.Sprintf("%s/%s", b.CurrentDir, "genframe"), []string{}, &os.ProcAttr{
		Sys: &syscall.SysProcAttr{
			Pgid: 1111,
		},
	})
	if err != nil {
		return err
	}
	b.genframeProc = genframeProc

	go func() {
		_, err := b.genframeProc.Wait()
		b.StoppedGenframe(err)
	}()
	return nil
}

func (b *Bootstrap) StoppedGenframe(err error) {
	b.StartGenframe()
}

func (b *Bootstrap) StopGenframe() {
	b.genframeProc.Kill()
}

func (b *Bootstrap) StartAutoupdater() error {
	autoupdaterProc, err := os.StartProcess(fmt.Sprintf("%s/%s", b.CurrentDir, "autoupdater"), []string{}, &os.ProcAttr{
		Sys: &syscall.SysProcAttr{
			Pgid: 1111,
		},
	})
	if err != nil {
		return err
	}
	b.autoupdaterProc = autoupdaterProc

	go func() {
		state, err := b.autoupdaterProc.Wait()
		b.StoppedAutoupdater(state, err)
	}()
	return nil
}

func (b *Bootstrap) StoppedAutoupdater(state *os.ProcessState, err error) {
	fmt.Println("stopped autoupdater", state, err, b.autoupdaterProc)
	b.StartAutoupdater()
	fmt.Println("restarted autoupdater", b.autoupdaterProc)
}

func (b *Bootstrap) StopAutoupdater() {
	b.autoupdaterProc.Kill()
}

func (b *Bootstrap) Run() {
	// if err := b.StartGenframe(); err != nil {
	// 	panic(err)
	// }
	if err := b.StartAutoupdater(); err != nil {
		panic(err)
	}
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Println("app - Run - signal: " + s.String())
	case err := <-b.notify:
		fmt.Println(fmt.Errorf("app - Run - b.notify: %w", err))
	}

	b.StopAutoupdater()
	b.StopGenframe()
}
