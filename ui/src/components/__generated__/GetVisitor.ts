

/* tslint:disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: GetVisitor
// ====================================================

export interface GetVisitor_visitor_bookmarks_user {
  name: string;
}

export interface GetVisitor_visitor_bookmarks_entry {
  id: string;
  title: string;
  url: string;
}

export interface GetVisitor_visitor_bookmarks {
  id: string;
  user: GetVisitor_visitor_bookmarks_user;
  entry: GetVisitor_visitor_bookmarks_entry;
  comment: string;
}

export interface GetVisitor_visitor {
  name: string;
  bookmarks: GetVisitor_visitor_bookmarks[];
}

export interface GetVisitor {
  visitor: GetVisitor_visitor;
}

/* tslint:disable */
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

//==============================================================
// END Enums and Input Objects
//==============================================================