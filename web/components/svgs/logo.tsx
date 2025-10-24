export default function LogoSvg(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      width="250"
      height="250"
      viewBox="0 0 250 250"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <circle cx="125" cy="125" r="125" fill="black" />
      <path d="M125 70L172.631 152.5H77.3686L125 70Z" fill="white" />
    </svg>
  );
}
