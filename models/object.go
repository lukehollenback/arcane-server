package models

type Object interface {

  //
  // Provides the unique identifier of the object for use in things like lookup tables.
  //
  ObjectID() string

  //
  // Provides the identifier of the area that the object resides in.
  //
  AreaID() string

  //
  // Provides the "x" coordinate of the object.
  //
  X() int

  //
  // Provides the "y" coordinate of the object.
  //
  Y() int

  //
  // Provides the "z" coordinate of the object.
  //
  Depth() int

}