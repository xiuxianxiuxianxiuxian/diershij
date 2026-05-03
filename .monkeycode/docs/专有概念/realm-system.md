# 境界系统

## 概述

修仙世界包含 10 个修仙境界，从凡人到渡劫飞升。境界是实体的核心属性，决定了寿命、天劫强度和基础能力。

## 境界等级

| 序号 | 常量名 | 值 | 中文 | 寿命 (年) | 天劫基础概率 | 天劫类型 | 强度倍率 |
|------|--------|-----|------|-----------|-------------|----------|----------|
| 1 | RealmMortal | mortal | 凡人 | 80 | N/A | N/A | N/A |
| 2 | RealmQiCondensation | qi_condensation | 凝气期 | 120 | 0.1 | thunder (雷劫) | 1.0x |
| 3 | RealmFoundation | foundation | 筑基期 | 200 | 0.2 | thunder (雷劫) | 2.0x |
| 4 | RealmGoldenCore | golden_core | 金丹期 | 500 | 0.3 | thunder_fire (雷火劫) | 5.0x |
| 5 | RealmNascentSoul | nascent_soul | 元婴期 | 1000 | 0.5 | thunder_fire_wind (雷火风劫) | 10.0x |
| 6 | RealmSoulTransform | soul_transformation | 化神期 | 3000 | 0.6 | five_element (五行劫) | 20.0x |
| 7 | RealmVoidRefinement | void_refinement | 炼虚期 | 5000 | 0.7 | heart_demon (心魔劫) | 50.0x |
| 8 | RealmIntegration | integration | 合体期 | 8000 | 0.8 | dao_tribulation (道劫) | 100.0x |
| 9 | RealmMahayana | mahayana | 大乘期 | 10000 | 0.9 | extinction (寂灭劫) | 200.0x |
| 10 | RealmTribulation | tribulation | 渡劫期 | 15000 | 0.95 | ascension (飞升劫) | 500.0x |

## 突破机制

### 前置条件

- `CultivationProgress >= 100` (修炼进度满)

### 成功率计算

```
successRate = 0.5 + (Luck / 200)
上限: 0.8 (80%)
```

运气值越高，突破成功率越高。基础 50%，运气满分 (100) 时达到上限 80%。

### 突破效果

成功突破后：
1. 境界提升到下一级
2. 修炼进度重置为 0
3. `MaxQi *= 1.5` (最大灵气提升 50%)
4. `MaxSpiritualPower *= 1.5` (最大灵力提升 50%)
5. `MaxLifespan` 更新为新境界对应的寿命

### 突破失败

当前代码中突破操作没有失败处理逻辑（成功率为计算值但实际判断逻辑简化），待确认失败时的具体行为。

## 修炼机制

### 修炼进度计算

```
cultivationGain = 0.1 * (Comprehension / 100)
```

悟性 (Comprehension) 越高，每次修炼获得的修为越多。

### 修炼状态

修炼时实体状态变为 `StatusCultivating`。

## 境界与属性

不同境界的实体在创建时具有不同的基础属性：

**凡人初始值**:
- Qi: 100, MaxQi: 100
- SpiritualPower: 100, MaxSpiritualPower: 100
- DivineSense: 10
- Comprehension: 50
- Constitution: 50
- Luck: 50
- AttackPower: 10
- Defense: 10
- Speed: 10
- MentalStability: 50
- RemainingLifespan / MaxLifespan: 80

## 天道印记

突破时不会直接改变天道印记。天道印记由因果系统 (`heavenly-dao` 服务) 独立管理，根据实体的行为善恶进行评估。
