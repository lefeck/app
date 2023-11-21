package main

import (
	"fmt"
	"log"
	"time"

	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"github.com/unidoc/unioffice/spreadsheet"
)

type biaozhi struct {
	yeMeiTime        string // 页眉的时间编号
	neiRouTime       string // 内容的时间格式
	ipaddr           string // 其实就是告警对象：主要填IP
	tiaomushu        string // 告警的条目数
	gaojingxinxi     string // 告警信息：告警描述
	xitong           string // 被攻击的系统名称
	danwei           string // 被攻击的单位
	chuzhijianyi     string // 处置建议
	lujing           string // 告警中的路径，病毒存在的路径、webshell的路径等
	shijianmingcheng string // 事件名称
	yjShiJianLeiXing string // 应急事件报告中特有的事件类型
	shijianjibie     string // 事件的级别
	loudongmingcheng string // 漏洞预警：漏洞名称
	yxBanBen         string // 漏洞预警：影响版本
	wenjianmingcheng string // 另存为文件的名称
	fengxianmiaoshu  string // 风险验证描述
}

var (
	// Baogao 要打开的报告类型
	Baogao string
	// SaveFile 要打开的报告类型
	SaveFile string
	dataList biaozhi
)

// 1.判断输出的是什么通告
func tonggao() {
	a, _ := spreadsheet.Open("素材.xlsx")
	excel, err := a.GetSheet("需要填写的资料")
	if err != nil {
		log.Fatalf("excel打开表错误，错误信息: %s", err)
	}
	defer a.Close()

	SaveFile = excel.Cell("B1").GetString() // 看是什么类型的报告获取字符串。
	switch SaveFile {
	case "可疑行为":
		Baogao = "可疑行为预警通告.docx"
	case "风险预警":
		Baogao = "安全风险预警通告.docx"
	case "应急预警":
		Baogao = "安全事件应急通告.docx"
	case "漏洞预警":
		Baogao = "安全漏洞预警通告.docx"
	default:
		log.Println("请检查(素材.xlsx)文件（需要填写的资料）表中报告类型是否正确！")
	}

	// 获取相关的数据
	dataList.shijianmingcheng = excel.Cell("B2").GetString()  // 风险预警：获取事件名称
	dataList.shijianjibie = excel.Cell("B3").GetString()      // 所有预警：获取事件级别
	dataList.ipaddr = excel.Cell("B4").GetString()            // 应急预警：业务IP地址
	dataList.lujing = excel.Cell("B5").GetString()            // 获取告警路径；
	dataList.tiaomushu = excel.Cell("B6").GetString()         // 获取告警路径；
	dataList.gaojingxinxi = excel.Cell("B7").GetString()      // 所有预警：告警描述
	dataList.chuzhijianyi = excel.Cell("B8").GetString()      // 所有预警：获取处置建议
	dataList.yjShiJianLeiXing = excel.Cell("B9").GetString()  // 应急预警：事件类型
	dataList.loudongmingcheng = excel.Cell("B10").GetString() // 漏洞预警：漏洞名称
	dataList.yxBanBen = excel.Cell("B11").GetString()         // 漏洞预警：影响版本
	dataList.xitong = excel.Cell("B12").GetString()           // 风险预警：业务系统名称
	dataList.danwei = excel.Cell("B13").GetString()           // 风险预警：业务单位名称
	dataList.fengxianmiaoshu = excel.Cell("B14").GetString()  // 风险预警：风险验证描述
}

// 2.处理并输出报告
func baogao() {

	// 打开一个已有格式的文档
	doc, err := document.Open(Baogao)
	if err != nil {
		log.Fatalf("打开word文档错误，错误信息: %s", err)
	}

	// 设置页眉
	hdr := doc.AddHeader()
	para := hdr.AddParagraph()
	para.Properties().AddTabStop(8*measurement.Inch, wml.ST_TabJcRight, wml.ST_TabTlcUnset)
	run := para.AddRun()
	run.AddTab()
	run.AddText("编号:-" + dataList.yeMeiTime)
	run.Properties().SetBold(true)
	run.Properties().SetSize(8)
	// run.Properties().SetColor(color.Black)
	doc.BodySection().SetHeader(hdr, wml.ST_HdrFtrUnset)
	// 页眉设置完毕

	for _, p := range doc.Paragraphs() { // 遍历所有段落
		for _, r := range p.Runs() { // 遍历相同格式的段落
			if SaveFile == "应急预警" {
				switch r.Text() {
				case "事件名称：":
					r.ClearContent()
					r.AddText("事件名称：" + dataList.shijianmingcheng)
				case "事件类型：":
					r.ClearContent()
					r.AddText("事件类型：" + dataList.yjShiJianLeiXing)
				case "事件级别：":
					r.ClearContent()
					r.AddText("事件级别：" + dataList.shijianjibie)
				case "告警对象：":
					r.ClearContent()
					r.AddText("告警对象：" + dataList.ipaddr)
				case "所属系统：":
					r.ClearContent()
					r.AddText("所属系统：" + dataList.xitong + "系统,隶属于" + dataList.danwei)
				case "告警内容：":
					r.ClearContent()
					r.AddText("告警内容：" + dataList.gaojingxinxi)
				case "路径":
					r.ClearContent()
					r.AddText(dataList.lujing)
				case "处置建议哦":
					r.ClearContent()
					r.AddText(dataList.chuzhijianyi)
				case "应急描述":
					r.ClearContent()
					r.AddText(dataList.neiRouTime + "，监控平台动态智能监控发现业务(IP：" + dataList.ipaddr + ")产生" + dataList.tiaomushu + "条“" + dataList.shijianmingcheng + "”日志告警信息。经资产确认，该主机为“" + dataList.xitong + "系统,隶属于" + dataList.danwei + "。")
				}
			} else {
				switch r.Text() {
				case "漏洞名称哦":
					r.ClearContent()
					r.AddText(dataList.loudongmingcheng) // 漏洞名称
				case "影响版本哦":
					r.ClearContent()
					r.AddText(dataList.yxBanBen) // 影响版本
				case "风险描述哦":
					r.ClearContent()
					r.AddText(dataList.gaojingxinxi) // 风险描述
				case "事件名称哦":
					r.ClearContent()
					r.AddText(dataList.shijianmingcheng)
				case "风险等级哦":
					r.ClearContent()
					r.AddText(dataList.shijianjibie) //风险等级
				case "疑点1":
					r.ClearContent()
					r.AddText(dataList.gaojingxinxi)
				case "路径":
					r.ClearContent()
					r.AddText(dataList.lujing)
				case "处置建议哦":
					r.ClearContent()
					r.AddText(dataList.chuzhijianyi)
				case "业务所属":
					yewu := "经资产确认，该主机为" + dataList.xitong + "系统,隶属于" + dataList.danwei
					r.ClearContent()
					r.AddText(yewu)
				case "风险验证描述":
					r.ClearContent()
					r.AddText(dataList.fengxianmiaoshu)
				default:
					fmt.Println("not modifying", r.Text())
				}
			}
		}
	}

	switch SaveFile {
	case "可疑行为":
		dataList.wenjianmingcheng = "【" + dataList.shijianmingcheng + "】" + "可疑行为预警通告" + dataList.yeMeiTime + ".docx"
	case "风险预警":
		dataList.wenjianmingcheng = "【" + dataList.shijianmingcheng + "】" + "安全风险预警通告" + dataList.yeMeiTime + ".docx"
	case "应急预警":
		dataList.wenjianmingcheng = "【" + dataList.shijianmingcheng + "】" + "安全事件应急通告" + dataList.yeMeiTime + ".docx"
	case "漏洞预警":
		dataList.wenjianmingcheng = "【" + dataList.loudongmingcheng + "】" + "安全漏洞预警通告" + dataList.yeMeiTime + ".docx"
	}
	err = doc.SaveToFile(dataList.wenjianmingcheng)
	if err != nil {
		log.Fatalln("文件保存错误：", err)
	}
}

func main() {
	// 获取符合我们格式的时间。
	dataList.neiRouTime = time.Now().Format("2006年01月02日15时04分")
	dataList.yeMeiTime = time.Now().Format("200601021504")

	tonggao()

	baogao()

}
