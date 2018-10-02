import React from "react";
import {Mutation, MutationUpdaterFn} from "react-apollo";
import gql from "graphql-tag";

import {bookmarkListItemFragment} from "./bookmarks";
import {CreateBookmark, CreateBookmarkVariables} from "./__generated__/CreateBookmark";
import {query as getVisitorQuery} from "./me";
import {GetVisitor} from "./__generated__/GetVisitor";

const mutation = gql`mutation CreateBookmark($url: String!, $comment: String!) {
  createBookmark(url: $url, comment: $comment) {
    id
    ...BookmarkListItemFragment
  }
}
${bookmarkListItemFragment}
`;

const updateBookmarks: MutationUpdaterFn<CreateBookmark> = (cache, result) => {
  const visitor = cache.readQuery<GetVisitor>({ query: getVisitorQuery });
  const { data } = result;
  if (visitor && data) {
    const bookmarks = [...visitor.visitor.bookmarks];
    const found = bookmarks.findIndex(bookmark => bookmark.id === data.createBookmark.id);
    if (found !== -1) {
      bookmarks[found] = data.createBookmark;
    } else {
      bookmarks.unshift(data.createBookmark);
    }
    const newVisitor = {
      visitor: {
        ...visitor.visitor,
        bookmarks,
      }
    };
    cache.writeQuery({ query: getVisitorQuery, data: newVisitor });
  }
};

export const AddBookmark: React.StatelessComponent = () => (
  <div className="AddBookmark">
    <Mutation<CreateBookmark, CreateBookmarkVariables> mutation={mutation} update={updateBookmarks}>
      {(createBookmark) => {
        return <BookmarkForm create={(url: string, comment: string) => {
          createBookmark({ variables: { url, comment } })
        }} />;
      }}
    </Mutation>
  </div>
);


interface BookmarkFormProps {
  create: (url: string, comment: string) => void;
}
interface BookmarkFormState {
  url: string
  comment: string
}
class BookmarkForm extends React.PureComponent<BookmarkFormProps, BookmarkFormState> {
  state = {
    url: '',
    comment: '',
  };

  private handleInput = (event: React.ChangeEvent<HTMLInputElement>) => {
    const input = event.currentTarget;
    switch (input.name) {
      case "url":
        this.setState({
          url: input.value,
        });
        break;
      case "comment":
        this.setState({
          comment: input.value,
        });
        break;
    }
  };

  private handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    this.props.create(this.state.url, this.state.comment);
  };

  render() {
    return (
      <form className="BookmarkForm" onSubmit={this.handleSubmit}>
        <label>URL:
          <input type="url" name="url" value={this.state.url} onChange={this.handleInput} />
        </label>
        <label>Comment:
          <input type="text" name="comment" value={this.state.comment} onChange={this.handleInput} />
        </label>
        <button>Create</button>
      </form>
    );
  }
}
