export interface Entity {
  id: string
  entity_type: 'player' | 'npc'
  name: string
  realm: string
  position: {
    region_id: string
    x: number
    y: number
  }
  attributes: {
    qi: number
    max_qi: number
    spiritual_power: number
    max_spiritual_power: number
    divine_sense: number
    comprehension: number
    constitution: number
    luck: number
    cultivation_progress: number
    attack_power: number
    defense: number
    speed: number
    mental_stability: number
    remaining_lifespan: number
    max_lifespan: number
  }
  karma: {
    karma_value: number
    merit: number
    heavenly_mark: string
  }
  status: string
  created_at: string
  updated_at: string
}

export interface Message {
  type: string
  payload: Record<string, unknown>
  timestamp: number
  request_id?: string
}

export interface OperationResult {
  success: boolean
  message: string
  effects: Record<string, unknown>
  timestamp: number
}

export interface Region {
  id: string
  name: string
  parent_region_id?: string
  spiritual_density: number
  spiritual_tier: number
  danger_level: number
  description: string
  resources: Resource[]
}

export interface Resource {
  id: string
  name: string
  type: string
  rarity: number
  quantity: number
}

export interface ChatMessage {
  id: string
  sender_id: string
  sender_name: string
  channel: string
  content: string
  timestamp: number
}
