# ==========================================
# 修炼与战斗规则 (Cultivation & Combat)
# ==========================================

## 1. 修炼效率规则

### 1.1 基础修炼速率
```
rate = comprehension × 0.1 × spiritual × method_match × realm_penalty × mental × (1 - aging)
```

**参数说明**:
- `comprehension`: 悟性值 (1-100)
- `spiritual`: `location_spiritual_density / 100`
- `method_match`: 灵根与功法匹配度 (0-1)
  - `match_count / max(required_roots_count, 1)`
- `realm_penalty`: `1.0 / (1.0 + realm_level × 0.2)`
- `mental`: 心境系数，稳定度 > 80 为 1.0，低于 50 线性衰减
- `aging`: 衰老惩罚 (0-0.2)

### 1.2 灵气潮汐修正
- 涨潮期：修炼速率 × 1.3
- 退潮期：修炼速率 × 0.7
- 周期：365 天（正弦函数）

## 2. 功法冲突规则

### 2.1 冲突检测
同时修炼多个功法时，计算冲突值：

| 冲突类型 | 冲突值增加 |
|----------|------------|
| 属性相克 (如金木) | +0.5 |
| 阵营对立 (正 vs 魔) | +0.3 |
| 境界要求差异 > 2 级 | +0.2 |

```
conflict_score = sum(conflicts) / (n × (n-1) / 2)
```

### 2.2 反噬概率
| 冲突值范围 | 反噬概率 | 效果 |
|------------|----------|------|
| < 0.2 | 0% | 无 |
| 0.2 - 0.5 | 10% | 修炼速率 -20% |
| 0.5 - 0.8 | 30% | 修炼速率 -50%，可能受伤 |
| > 0.8 | 60% | 修炼速率 -80%，高概率重伤 |

## 3. 战斗伤害规则

### 3.1 伤害计算公式
```
damage = base_atk × skill_mult × realm_suppression × element × (1 - def_reduction) × method_bonus
```

**修正因子**:
- `realm_suppression = 1.0 + realm_diff × 0.15`
  - 每高 1 境界 +15%，低 1 境界 -15%
- `element`: 五行克制系数 (1.0-1.5)
- `def_reduction = defense / (defense + 100)`

### 3.2 五行克制表

| 攻击属性 | 克制 (+50%) | 被克 (-50%) |
|----------|-------------|-------------|
| 火 (Fire) | 金、木 | 水 |
| 水 (Water) | 火 | 土 |
| 木 (Wood) | 土 | 金 |
| 金 (Metal) | 木 | 火 |
| 土 (Earth) | 水 | 木 |

### 3.3 暴击规则
- 基础暴击率：5%
- 基础暴击伤害：150%
- 可通过功法/装备提升

### 3.4 战斗结算
- 气血归零 → 死亡或重伤
- 重伤 → 境界跌落概率 10%，物品掉落
- 死亡 → 境界跌落，极品灵石掉落 50%，因果业力增加

## 4. 法则与领域规则

### 4.1 法则感悟
- 基础感悟速率：0.1/天
- 功法亲和倍率：`law_comprehension_bonus`
- 灵气浓度修正：同修炼效率

### 4.2 领域展开条件
- 境界要求：化神及以上
- 法则感悟要求：至少一项法则 ≥ 50
- 领域强度：`sum(law_values) / 10`
- 领域范围：`domain_power × 10` 米
- 领域内压制：`law_suppression × 0.1` 倍敌方全属性

## 5. 心魔入侵规则

### 5.1 触发概率 (每日)
```
P = 0.001 × mental_factor × karma_factor × breakthrough_stress × conflict_factor
```

- `mental_factor = (100 - stability) / 100`
- `karma_factor = 1.0 + karma / 1000`
- `breakthrough_stress`: 7 天内尝试突破为 2.0
- `conflict_factor`: `1.0 + conflict_score × 2.0`

### 5.2 心魔战果
| 结果 | 效果 |
|------|------|
| 战胜 | 心境 +20，修为进度 +10% |
| 战败 | 修为进度 -20%，心境 -30 |
| 战败 + 走火入魔 (30%) | 境界跌落，属性永久 -10% |
