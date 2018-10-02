import React from "react";
import {Link} from "react-router-dom";
import gql from "graphql-tag";

import {BookmarkListFragment} from "./__generated__/BookmarkListFragment";
import {BookmarkListItemFragment} from "./__generated__/BookmarkListItemFragment";

export const bookmarkListItemFragment = gql`fragment BookmarkListItemFragment on Bookmark {
  user {
      name
  }
  entry {
    id
    title
    url
  }
  comment
}
`;

interface BookmarkListItemProps {
  bookmark: BookmarkListItemFragment
  deleteBookmark?: () => void
}
const BookmarkListItem: React.StatelessComponent<BookmarkListItemProps> = ({ bookmark, deleteBookmark }) => (
  <div className="BookmarkListItem">
    <div>
      <Link to={`/entry/${bookmark.entry.id}`}>{bookmark.entry.title}</Link>
      <span> - </span>
      <a href={bookmark.entry.url}>{bookmark.entry.url}</a>
    </div>
    <div>
      <h2>{bookmark.user.name}</h2>
      <p>{bookmark.comment}</p>
      {deleteBookmark && <button onClick={deleteBookmark}>Delete</button>}
    </div>
  </div>
);

export const bookmarkListFragment = gql`fragment BookmarkListFragment on User {
  name
  bookmarks {
    id
    ...BookmarkListItemFragment
  }
}
${bookmarkListItemFragment}
`;

interface BookmarkListProps {
  user: BookmarkListFragment
  deleteBookmark?: (bookmarkId: string) => void
}
export const BookmarkList: React.StatelessComponent<BookmarkListProps> = ({ user, deleteBookmark }) => (
  <div className="BookmarkList">
    <h1>{user.name}‘s bookmarks</h1>
    <ul>
      {user.bookmarks.map(bookmark => (<li key={bookmark.id}>
        <BookmarkListItem bookmark={bookmark} deleteBookmark={deleteBookmark ? () => { deleteBookmark(bookmark.id) } : undefined} />
      </li>))}
    </ul>
  </div>
);
