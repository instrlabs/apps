import React from "react";
import clsx from "clsx";

interface FormInputProps {
  id: string;
  type: string;
  label: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder: string;
  isInvalid?: boolean;
  errorMessage?: string;
}

const FormInput: React.FC<FormInputProps> = ({
  id,
  type,
  value,
  onChange,
  placeholder,
  isInvalid,
  errorMessage,
}) => (
  <div className="relative">
    <input
      id={id}
      type={type}
      placeholder={placeholder}
      className={clsx(
        "w-full px-5 py-4 bg-white shadow-primary rounded-xl",
        isInvalid ? "outline outline-red-300 focus:outline focus:outline-red-300" : "focus:outline-none",
      )}
      value={value}
      onChange={onChange}
      aria-invalid={isInvalid ? "true" : "false"}
      aria-errormessage={isInvalid && errorMessage ? `${id}-error` : undefined}
      aria-describedby={isInvalid && errorMessage ? `${id}-error` : undefined}
      required
    />
    {isInvalid && errorMessage ? (
      <span
        id={`${id}-error`}
        role="alert"
        className="absolute left-5 top-full mt-1 text-xs text-red-400"
      >
        {errorMessage}
      </span>
    ) : null}
  </div>
);

export default FormInput;
