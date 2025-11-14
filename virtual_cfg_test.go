package valuerange

////////////////////////////////////////////////////////////////////////////////
/// 这边声明一些，模拟策划配置的数据，方便其他 test 使用

type heroTagCfg struct {
	Free bool // 免费使用英雄
}

type heroCfg struct {
	Id      uint64            // 配置id
	Desc    string            // 描述
	Quality int               // 品质
	Open    bool              // 是否开放可以使用
	Tag     heroTagCfg        // 英雄标签配置
	Skins   []uint64          // 可以使用的皮肤配置id列表
	Attrs   map[uint32]uint32 // 英雄属性列表 key: 属性id value: 属性值
}

// 假设有这么一张 hero.csv 配置表的数据
const heroCfgKey = "heroCfg"

var heroCfgList = []heroCfg{
	{Id: 61401, Desc: "top", Quality: 1, Open: true, Tag: heroTagCfg{Free: true}, Skins: []uint64{6140101, 6140102}, Attrs: map[uint32]uint32{1: 100, 2: 200}},
	{Id: 61402, Desc: "ace", Quality: 2, Open: true, Tag: heroTagCfg{Free: false}, Skins: []uint64{6140201, 6140202}, Attrs: map[uint32]uint32{1: 150, 2: 250}},
	{Id: 61403, Desc: "mid", Quality: 5, Open: false, Tag: heroTagCfg{Free: true}, Skins: []uint64{6140301, 6140302}, Attrs: map[uint32]uint32{1: 200, 2: 300}},
}

type heroSkinCfg struct {
	Id   uint64 // 皮肤配置id
	Desc string // 皮肤描述
}

// 假设有这么一张 heroSkin.csv 配置表的数据
const heroSkinCfgKey = "heroSkinCfg"

var heroSkinCfgList = []heroSkinCfg{
	{Id: 6140101, Desc: "top skin 1"},
	{Id: 6140102, Desc: "top skin 2"},
	{Id: 6140201, Desc: "ace skin 1"},
	{Id: 6140202, Desc: "ace skin 2"},
	{Id: 6140301, Desc: "mid skin 1"},
	{Id: 6140302, Desc: "mid skin 2"},
}
