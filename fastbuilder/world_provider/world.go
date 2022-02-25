package world_provider

import (
	"phoenixbuilder/dragonfly/server/world"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
)

var CurrentWorld *world.World = nil

// we let this function auto init when programme start
// not init it in every session, nor explicit init it, pay attention please
func init() {
	InitRuntimeIdsWithoutMinecraftPrefix()
}

func Create(conn *minecraft.Conn) *world.World {
	intw:=world.New(&StubLogger{},32)
	intw.Provider(NewOnlineWorldProvider(conn))
	return intw
}

func NewWorld(conn *minecraft.Conn) {
	ChunkCache=make(map[world.ChunkPos]*packet.LevelChunk)
	CurrentWorld=Create(conn)
	firstLoaded=false
}

func DestroyWorld() {
	firstLoaded=false
	CurrentWorld=nil
	ChunkCache=nil
}

