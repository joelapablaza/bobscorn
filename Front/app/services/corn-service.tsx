const uri = process.env.NEXT_PUBLIC_URI_SERVER;

export async function purchaseCornFromBob() {
  try {
    if (!uri) {
      throw new Error('NEXT_PUBLIC_URI_SERVER is not defined');
    }
    console.log(`Corn purchase request - URI: ${uri}`);
    const response = await fetch(uri, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({}),
    });

    const data = await response.json();

    console.log(
      `Corn purchase response - Status: ${response.status}, Message:`,
      data
    );

    if (response.ok) {
      return {
        success: true,
        message: data.message,
        status: response.status,
      };
    } else {
      return {
        success: false,
        error: data.error,
        status: response.status,
      };
    }
  } catch (networkError) {
    console.error('Network error connecting to corn server:', networkError);
    return {
      success: false,
      error: "Could not connect to Bob's farm. Is the server running?",
      status: 500,
    };
  }
}
