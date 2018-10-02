

/* tslint:disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL fragment: BookmarkListFragment
// ====================================================

export interface BookmarkListFragment_bookmarks_user {
  name: string;
}

export interface BookmarkListFragment_bookmarks_entry {
  id: string;
  title: string;
  url: string;
}

export interface BookmarkListFragment_bookmarks {
  id: string;
  user: BookmarkListFragment_bookmarks_user;
  entry: BookmarkListFragment_bookmarks_entry;
  comment: string;
}

export interface BookmarkListFragment {
  name: string;
  bookmarks: BookmarkListFragment_bookmarks[];
}

/* tslint:disable */
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

//==============================================================
// END Enums and Input Objects
//==============================================================