import React from "react";
import {Mutation, MutationUpdaterFn, Query} from "react-apollo";
import gql from "graphql-tag";

import {AddBookmark} from "./add";
import {BookmarkList, bookmarkListFragment} from "./bookmarks";
import {GetVisitor} from "./__generated__/GetVisitor";
import {DeleteBookmark, DeleteBookmarkVariables} from "./__generated__/DeleteBookmark";
import {CreateBookmark} from "./__generated__/CreateBookmark";

export const query = gql`query GetVisitor {
  visitor {
    ...BookmarkListFragment
  }
}
${bookmarkListFragment}
`;

export const deleteBookmark = gql`mutation DeleteBookmark($bookmarkId: ID!) {
  deleteBookmark(bookmarkId: $bookmarkId)
}`;

const updateBookmarks: (bookmarkId: string) => MutationUpdaterFn<DeleteBookmark> = (bookmarkId) => (cache, result) => {
  const visitor = cache.readQuery<GetVisitor>({ query });
  const { data } = result;
  if (visitor && data) {
    const bookmarks = [...visitor.visitor.bookmarks].filter(bookmark => bookmark.id !== bookmarkId);
    const newVisitor = {
      visitor: {
        ...visitor.visitor,
        bookmarks,
      }
    };
    cache.writeQuery({query, data: newVisitor});
  }
};

export const Me: React.StatelessComponent = () => (
  <div className="Me">
    <Query<GetVisitor> query={query}>
      {result => {
        if (result.error) {
          return <p className="error">Error: {result.error.message}</p>
        }
        if (result.loading) {
          return <p className="loading">Loading</p>
        }
        const {data} = result;
        return <>
          <AddBookmark />
          <Mutation<DeleteBookmark, DeleteBookmarkVariables> mutation={deleteBookmark}>
            {(deleteBookmark) => {
              return <BookmarkList
                user={data!.visitor}
                deleteBookmark={(bookmarkId: string) =>
                  deleteBookmark({ variables: { bookmarkId }, update: updateBookmarks(bookmarkId) })} />
            }}
          </Mutation>
        </>;
      }}
    </Query>
  </div>
);
