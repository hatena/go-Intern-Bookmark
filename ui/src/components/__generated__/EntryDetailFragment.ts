

/* tslint:disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL fragment: EntryDetailFragment
// ====================================================

export interface EntryDetailFragment_bookmarks_user {
  name: string;
}

export interface EntryDetailFragment_bookmarks {
  id: string;
  user: EntryDetailFragment_bookmarks_user;
  comment: string;
}

export interface EntryDetailFragment {
  title: string;
  url: string;
  bookmarks: EntryDetailFragment_bookmarks[];
}

/* tslint:disable */
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

//==============================================================
// END Enums and Input Objects
//==============================================================