package msgmodels

//
// ObjSync represents the structure of a message that tells clients to synchronize variables of an
// object amongst connected clients.
//
type ObjSync struct {
  ObjectID  string                 // The unique ID of the object instance. For use in things like synchronized object lookup tables.
  AreaID    string                 // The unique ID of the area that the object exists in. Helps with message broadcast localization.
  Variables map[string]interface{} // The actual payload of variables to synchronize.
}
