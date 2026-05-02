# ==========================================
# 炼丹与炼器规则 (Alchemy & Forging)
# ==========================================

## 1. 炼丹规则

### 1.1 成功率计算
```
P = base_rate × skill_factor × control × material × (1 + furnace_bonus) × spiritual
```

**参数说明**:
- `skill_factor = alchemy_level / recipe_required_level`
- `control = 0.5 + fire_control × 0.5` (火候控制)
- `material = avg(materials.quality) / 100`
- `furnace_bonus`: 丹炉品质加成
- `spiritual`: `location_spiritual_density / 100`

最终概率限制在 `[5%, 95%]` 区间。

### 1.2 品质生成
基于成功率的随机判定：

| 品质 | 概率区间 | 效果倍率 |
|------|----------|----------|
| 极品 | `roll < P × 0.1` | 3.0x |
| 上品 | `P × 0.1 ≤ roll < P × 0.3` | 2.0x |
| 中品 | `P × 0.3 ≤ roll < P` | 1.0x |
| 下品 (废丹) | `roll ≥ P` | 0.3x，可能带毒性 |

### 1.3 失败后果
| 结果 | 概率 | 效果 |
|------|------|------|
| 炸炉 | 20% | 气血 -30%，材料全损，可能引起注意 |
| 材料损毁 | 50% | 材料全损，无其他影响 |
| 废丹 | 30% | 产出下品丹药，效果差且有副作用 |

### 1.4 丹药品阶体系
- 一品至九品，每品分下/中/上/极品
- 品阶越高，所需材料越稀有，成功率越低
- 极品丹药可引发天劫或天地异象

## 2. 炼器规则

### 2.1 成功率计算
```
P = base_rate × skill_factor × material_match × fire_quality × (1 + formation_bonus)
```

**参数说明**:
- `skill_factor = artificing_level / blueprint_required_level`
- `material_match`: 材料兼容度 (0-1)
- `fire_quality = flame_tier / 10`
- `formation_bonus`: 辅助阵法数量 × 0.05

最终概率限制在 `[5%, 90%]` 区间。

### 2.2 法宝品阶
| 品阶 | 概率区间 | 特性 |
|------|----------|------|
| 古宝 | `roll < P × 0.05` | 自带神通，可成长 |
| 天器 | `P × 0.05 ≤ roll < P × 0.15` | 蕴含法则碎片 |
| 地器 | `P × 0.15 ≤ roll < P × 0.4` | 属性大幅增强 |
| 凡器 | `P × 0.4 ≤ roll < P` | 基础属性增强 |
| 失败 | `roll ≥ P` | 材料损毁 |

### 2.3 法宝绑定
- 法宝可认主，认主后他人无法使用
- 解除绑定需特殊仪式或大能出手
- 主人死亡，法宝可能认新主或灵性消散

## 3. 阵法规则

### 3.1 阵法威力
```
power = base_power × skill × eye_quality × flag_quality × (1 + terrain_bonus)
```

- `eye_quality`: 阵眼物品品质 / 100
- `flag_quality`: 阵旗平均品质 / 100
- `terrain_bonus`: 地势加成 (0-0.5)

### 3.2 破阵规则
```
P_break = breaker_power / (formation_power + breaker_power) × (1 + realm_diff × 0.1)
```

- 破阵失败可能触发阵法反击
- 高级阵法破阵失败有生命危险

## 4. 灵根觉醒规则

### 4.1 觉醒概率 (每日检测)
```
P = 0.001 × (1 + cultivation_years × 0.1) × (1 + breakthrough_count × 0.2) × (1 + luck / 50)
```

### 4.2 变异灵根
某些灵根组合可能产生变异：

| 原始组合 | 变异结果 | 特性 |
|----------|----------|------|
| 火 + 水 | 雷灵根 | 爆发力强，天劫加成 |
| 火 + 金 | 风灵根 | 速度极快，闪避高 |
| 水 + 土 | 冰灵根 | 控制力强，减速 |
| 木 + 火 | 光灵根 | 治疗加成，克魔 |
| 金 + 土 | 暗灵根 | 隐匿强，刺杀加成 |
