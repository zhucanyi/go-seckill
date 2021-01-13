package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhucanyi/go_seckill/common"
	"github.com/zhucanyi/go_seckill/jd_seckill"
	"github.com/zhucanyi/go_seckill/log"
	"os"
	"time"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open JD’s simulated login",
	Run:   startLogin,
}

func startLogin(cmd *cobra.Command, args []string) {
	session := jd_seckill.NewSession(common.CookieJar)
	//检测是否登录过
	if common.Exists("./cookie.txt") {
		//已登录，检测登录状态
		err := session.CheckLoginStatus()
		if err != nil {
			log.Error("登录失效，请重新登录")
			return
		}
		user := jd_seckill.NewUser(common.Client, common.Config)
		log.Info("登录成功")
		userInfo, _ := user.GetUserInfo()
		log.Info("用户:" + userInfo)
	} else {
		//未登录
		user := jd_seckill.NewUser(common.Client, common.Config)
		wlfstkSmdl, err := user.QrLogin()
		defer user.DelQrCode()
		if err != nil {
			os.Exit(0)
		}
		ticket := ""
		for {
			ticket, err = user.QrcodeTicket(wlfstkSmdl)
			if err == nil && ticket != "" {
				break
			}
			time.Sleep(2 * time.Second)
		}
		_, err = user.TicketInfo(ticket)
		if err == nil {
			if status := user.RefreshStatus(); status == nil {
				//保存cookie
				_ = session.SaveCookieToFile("./cookie.txt")
				log.Info("登录成功")
				userInfo, _ := user.GetUserInfo()
				log.Info("用户:" + userInfo)
			} else {
				log.Error("登录失效")
			}
		} else {
			log.Error("登录失败")
		}
	}
}
