

/* tslint:disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: GetEntry
// ====================================================

export interface GetEntry_getEntry_bookmarks_user {
  name: string;
}

export interface GetEntry_getEntry_bookmarks {
  id: string;
  user: GetEntry_getEntry_bookmarks_user;
  comment: string;
}

export interface GetEntry_getEntry {
  id: string;
  title: string;
  url: string;
  bookmarks: GetEntry_getEntry_bookmarks[];
}

export interface GetEntry {
  getEntry: GetEntry_getEntry;
}

export interface GetEntryVariables {
  entryId: string;
}

/* tslint:disable */
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

//==============================================================
// END Enums and Input Objects
//==============================================================