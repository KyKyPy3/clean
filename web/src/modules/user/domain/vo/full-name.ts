import { ValueObject } from '@base';
import { isEmpty } from '@utils';
import { ArgumentNotProvidedException } from '@core';

export interface FullNameProps {
  name: string;
  surname: string;
  middlename: string;
}

export class FullName extends ValueObject<FullNameProps> {
  get name(): string {
    return this.props.name
  }

  get surname(): string {
    return this.props.surname
  }

  get middlename(): string {
    return this.props.middlename
  }

  protected validate(props: FullNameProps): void {
    if (isEmpty(props.name)) {
      throw new ArgumentNotProvidedException("missing name in fullname");
    }
  }
}