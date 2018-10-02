

/* tslint:disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: CreateBookmark
// ====================================================

export interface CreateBookmark_createBookmark_user {
  name: string;
}

export interface CreateBookmark_createBookmark_entry {
  id: string;
  title: string;
  url: string;
}

export interface CreateBookmark_createBookmark {
  id: string;
  user: CreateBookmark_createBookmark_user;
  entry: CreateBookmark_createBookmark_entry;
  comment: string;
}

export interface CreateBookmark {
  createBookmark: CreateBookmark_createBookmark;
}

export interface CreateBookmarkVariables {
  url: string;
  comment: string;
}

/* tslint:disable */
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

//==============================================================
// END Enums and Input Objects
//==============================================================