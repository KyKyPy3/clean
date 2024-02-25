import { ValueObject } from "@base"
import { ArgumentNotProvidedException } from "@core"
import { isEmpty } from "@utils"

export interface EmailProps {
	value: string;
}

export class Email extends ValueObject<EmailProps> {
	get value(): string {
    return this.props.value
  }

	protected validate(props: EmailProps): void {
    if (isEmpty(props.value)) {
      throw new ArgumentNotProvidedException("missing name in fullname");
    }
  }
}
