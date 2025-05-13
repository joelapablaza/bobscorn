async function purchaseCornFromBob() {
  try {
    const response = await fetch('http://localhost:8000/buy', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    } else {
      const text = await response.text();
      return {
        error: `Server error: ${text || response.statusText}`,
      };
    }
  } catch (error) {
    console.error('Error purchasing corn:', error);
    return { error: 'Failed to purchase corn. Please try again later.' };
  }
}
