package msgmodels

//
// ObjCreate represents the structure of a message that tells clients to create a new instance of a
// synchronized object.
//
type ObjCreate struct {
  Type     string // The name of the object to create. Should reference a literal object name/identifier on the client.
  ObjectID string // The unique ID of the object instance. For use in things like synchronized object lookup tables.
  AreaID   string // The unique ID of the area that the object should be created in.
  X        int    // The initial horizontal location of the relevant character instance.
  Y        int    // The initial vertical location of the relevant character instance.
  Depth    int    // The depth that the relevant character should exist in the room that that it is instantiated into.
}
