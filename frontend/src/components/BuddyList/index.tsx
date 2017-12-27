import * as React from 'react';
import Buddy from '../../models/Buddy';
import './index.css';

interface Props {
    users: Buddy[];
}
class BuddyList extends React.Component<Props> {
    public render() {
        return (
            <div className="buddy-list">
            <ul className="inner">
              {
                this.props.users.map((user: Buddy) => (
                  <li>{user.name}</li>
                ))
              }
            </ul>
          </div>
        );
    }
}

export default BuddyList;