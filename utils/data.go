package utils

import (
	"image/color"
)

var POTS_MAP = map[string]color.RGBA{"blindness": {64, 64, 64, 255}, "damageBoost": {139, 0, 0, 255}, "fireResistance": {255, 200, 0, 255}, "harm": {139, 0, 0, 255}, "heal": {255, 0, 0, 255}, "invisibility": {204, 204, 204, 255}, "jump": {0, 255, 0, 255}, "luck": {0, 255, 0, 255}, "moveSlowdown": {192, 192, 192, 255}, "moveSpeed": {0, 255, 255, 255}, "nightVision": {0, 0, 255, 255}, "poison": {0, 255, 0, 255}, "regeneration": {255, 175, 175, 255}, "slowFall": {192, 192, 192, 255}, "turtleMaster": {192, 192, 192, 255}, "waterBreathing": {192, 192, 192, 255}, "weakness": {0, 0, 0, 255}, "haste": {255, 255, 0, 255}, "wither": {61, 43, 31, 255}}

var ARMOR_MAP = map[string]string{"diamond_layer_1": "diamond_1", "diamond_layer_2": "diamond_2", "chainmail_layer_1": "chain_1", "chainmail_layer_2": "chain_2", "gold_layer_1": "gold_1", "gold_layer_2": "gold_2", "iron_layer_1": "iron_1", "iron_layer_2": "iron_2", "leather_layer_1": "cloth_1", "leather_layer_2": "cloth_2"}

var DEFAULT_RECOLOR_LIST = []string{"diamond_sword", "diamond_shovel", "diamond_pickaxe", "diamond_leggings", "diamond_horse_armor", "diamond_hoe", "diamond_chestplate", "diamond_helmet", "diamond_boots", "diamond_axe", "diamond", "diamond_block", "apple", "apple_golden", "bed", "icons", "widgets", "cubemap_0", "cubemap_1", "cubemap_2", "cubemap_3", "cubemap_4", "cubemap_5", "sky", "cloud1", "inventory", "gui", "hotdogempty", "hotdogfull", "diamond_1", "diamond_2", "experiencebarfull", "filled_progress_bar", "empty_progress_bar", "diamond_layer_1", "diamond_layer_2", "diamond_1", "diamond_2", "ender_pearl", "bow", "bow_pulling_0", "bow_pulling_1", "bow_pulling_2", "bow_standby", "totem", "generic_54", "pack_icon", "pack", "bed", "beacon", "bed_feet_end", "bed_feet_side", "bed_feet_top", "bed_head_end", "bed_head_side", "bed_head_top"}

var BEDS_MAP = map[string]color.RGBA{"black": {20, 20, 20, 255}, "blue": {20, 20, 255, 255}, "brown": {120, 75, 40, 255}, "cyan": {20, 165, 165, 255}, "gray": {120, 120, 120, 255}, "green": {100, 200, 20, 255}, "light_blue": {40, 150, 220, 255}, "lime": {120, 255, 25, 255}, "magenta": {190, 60, 180, 255}, "orange": {255, 127, 25, 255}, "pink": {245, 140, 170, 255}, "purple": {130, 45, 180, 255}, "silver": {180, 180, 180, 255}, "white": {240, 240, 240, 255}, "yellow": {255, 190, 25, 255}}
