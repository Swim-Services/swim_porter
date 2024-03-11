package utils

import (
	"image/color"
)

var POTS_MAP = map[string]color.RGBA{"blindness": {102, 102, 102, 255}, "damageBoost": {139, 0, 0, 255}, "fireResistance": {255, 102, 0, 255}, "harm": {139, 0, 0, 255}, "heal": {255, 0, 0, 255}, "invisibility": {204, 204, 204, 255}, "jump": {0, 204, 0, 255}, "luck": {0, 204, 0, 255}, "moveSlowdown": {204, 204, 204, 255}, "moveSpeed": {51, 204, 255, 255}, "nightVision": {0, 0, 255, 255}, "poison": {0, 204, 0, 255}, "regeneration": {255, 192, 203, 255}, "slowFall": {204, 204, 204, 255}, "turtleMaster": {204, 204, 204, 255}, "waterBreathing": {204, 204, 204, 255}, "weakness": {0, 0, 0, 255}, "haste": {255, 255, 0, 255}, "wither": {61, 43, 31, 255}}

var ARMOR_MAP = map[string]string{"diamond_layer_1": "diamond_1", "diamond_layer_2": "diamond_2", "chainmail_layer_1": "chain_1", "chainmail_layer_2": "chain_2", "gold_layer_1": "gold_1", "gold_layer_2": "gold_2", "iron_layer_1": "iron_1", "iron_layer_2": "iron_2", "leather_layer_1": "cloth_1", "leather_layer_2": "cloth_2"}
