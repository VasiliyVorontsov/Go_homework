package main

import (
	"strings"
)

type Item struct {
	Name string
}

type Player struct {
	CurrentRoom *Room
	HasBackpack bool
	Inventory   map[string]*Item
}

type Room struct {
	Name        string
	Description string
	Items       map[string]*Item
	Exits       map[string]*Room
	Locked      bool
	LockMessage string
}

type World struct {
	Player *Player
	Rooms  map[string]*Room
}

var world *World

func initGame() {
	world = &World{}

	kitchen := &Room{
		Name:        "кухня",
		Description: "кухня, ничего интересного.",
		Items:       make(map[string]*Item),
		Exits:       make(map[string]*Room),
	}

	corridor := &Room{
		Name:        "коридор",
		Description: "ничего интересного.",
		Items:       make(map[string]*Item),
		Exits:       make(map[string]*Room),
	}

	room := &Room{
		Name:        "комната",
		Description: "ты в своей комнате.",
		Items:       make(map[string]*Item),
		Exits:       make(map[string]*Room),
	}

	street := &Room{
		Name:        "улица",
		Description: "на улице весна.",
		Items:       make(map[string]*Item),
		Exits:       make(map[string]*Room),
		Locked:      true,
		LockMessage: "дверь закрыта",
	}

	kitchen.Items["чай"] = &Item{Name: "чай"}
	room.Items["ключи"] = &Item{Name: "ключи"}
	room.Items["конспекты"] = &Item{Name: "конспекты"}
	room.Items["рюкзак"] = &Item{Name: "рюкзак"}

	kitchen.Exits["коридор"] = corridor
	corridor.Exits["кухня"] = kitchen
	corridor.Exits["комната"] = room
	corridor.Exits["улица"] = street
	room.Exits["коридор"] = corridor
	street.Exits["домой"] = corridor

	world.Player = &Player{
		CurrentRoom: kitchen,
		HasBackpack: false,
		Inventory:   make(map[string]*Item),
	}

	world.Rooms = map[string]*Room{
		"кухня":   kitchen,
		"коридор": corridor,
		"комната": room,
		"улица":   street,
	}
}

func handleCommand(command string) string {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return "неизвестная команда"
	}

	args := parts[1:]

	switch parts[0] {
	case "осмотреться":
		return handleLook()
	case "идти":
		if len(args) < 1 {
			return "куда идти?"
		}
		return handleGo(args[0])
	case "взять":
		if len(args) < 1 {
			return "что взять?"
		}
		return handleTake(args[0])
	case "надеть":
		if len(args) < 1 {
			return "что надеть?"
		}
		return handleWear(args[0])
	case "применить":
		if len(args) < 2 {
			return "применить что и к чему?"
		}
		return handleUse(args[0], args[1])
	default:
		return "неизвестная команда"
	}
}

func handleLook() string {
	p := world.Player
	room := p.CurrentRoom

	if room.Name == "кухня" {
		desc := "ты находишься на кухне, на столе: чай"
		if !p.HasBackpack {
			desc += ", надо собрать рюкзак и идти в универ."
		} else {
			desc += ", надо идти в универ."
		}
		desc += " можно пройти - " + getAvailableExits(room)
		return desc
	}

	if room.Name == "комната" {
		items := []string{}
		if room.Items["ключи"] != nil {
			items = append(items, "ключи")
		}
		if room.Items["конспекты"] != nil {
			items = append(items, "конспекты")
		}

		desc := ""
		if len(items) > 0 {
			desc += "на столе: " + strings.Join(items, ", ")
		}

		if room.Items["рюкзак"] != nil {
			if len(items) > 0 {
				desc += ", на стуле: рюкзак"
			} else {
				desc += "на стуле: рюкзак"
			}
		}

		if len(items) == 0 && room.Items["рюкзак"] == nil {
			desc += "пустая комната"
		}

		desc += ". можно пройти - " + getAvailableExits(room)
		return desc
	}

	return room.Description + " можно пройти - " + getAvailableExits(room)
}

func handleGo(direction string) string {
	p := world.Player
	room := p.CurrentRoom

	targetRoom, exists := room.Exits[direction]
	if !exists {
		return "нет пути в " + direction
	}

	if targetRoom.Locked {
		return targetRoom.LockMessage
	}

	world.Player.CurrentRoom = targetRoom
	switch targetRoom.Name {
	case "комната":
		return "ты в своей комнате. можно пройти - коридор"
	case "кухня":
		return "кухня, ничего интересного. можно пройти - коридор"
	case "коридор":
		return "ничего интересного. можно пройти - кухня, комната, улица"
	case "улица":
		return "на улице весна. можно пройти - домой"
	default:
		return handleLook()
	}
}

func handleTake(itemName string) string {
	p := world.Player
	room := p.CurrentRoom

	item, exists := room.Items[itemName]
	if !exists || item == nil {
		return "нет такого"
	}

	if itemName == "рюкзак" {
		return handleWear("рюкзак")
	}

	if !p.HasBackpack {
		return "некуда класть"
	}

	delete(room.Items, itemName)
	p.Inventory[itemName] = item
	return "предмет добавлен в инвентарь: " + itemName
}

func handleWear(itemName string) string {
	if itemName != "рюкзак" {
		return "неизвестная команда"
	}

	p := world.Player
	room := p.CurrentRoom

	item, exists := room.Items[itemName]
	if !exists || item == nil {
		return "нет такого"
	}

	delete(room.Items, itemName)
	p.HasBackpack = true
	return "вы надели: рюкзак"
}

func handleUse(itemName, target string) string {
	p := world.Player

	item, exists := p.Inventory[itemName]
	if !exists || item == nil {
		return "нет предмета в инвентаре - " + itemName
	}

	if itemName == "ключи" && target == "дверь" && p.CurrentRoom.Name == "коридор" {
		street := world.Rooms["улица"]
		street.Locked = false
		return "дверь открыта"
	}

	return "не к чему применить"
}

func getAvailableExits(room *Room) string {
	exits := make([]string, 0, len(room.Exits))
	if room.Name == "коридор" {
		exits = append(exits, "кухня", "комната", "улица")
	} else {
		for exit := range room.Exits {
			exits = append(exits, exit)
		}
	}
	return strings.Join(exits, ", ")
}
