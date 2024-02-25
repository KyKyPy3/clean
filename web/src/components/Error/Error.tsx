const Error: React.FC = () => {
  return(
    <div className='w-screen min-h-screen p-2 overflow-y-scroll flex flex-col justify-center items-center default-bg'>
      <h1 className='text-slate-300 text-9xl mb-5'>Oops!</h1>
      <h2 className='text-slate-300 text-3xl'>The requested resource was not found.</h2>
    </div>
  )
}

export default Error