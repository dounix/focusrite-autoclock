# Focusrite auto clock

If you use digital inputs to focusrite consumer(Scarlett, or anything that uses control server) products, you have to pick an internal clock or S/PDIF clock.  

Many S/PDIF sources(Fractal AxeFx or Kemper) must be the clock master, and will result in pops and clicks in recording if you use the internal clock

When set to S/PDIF clock source, if that S/PDIF source is turned off, you will lose all audio.

This is a minor annoyance, moreso if you use the focusrite as your primary interface.
The real issue is not remembering to change back to S/PDIF, which you will notice after you record something, and only notice the pops and clicks when you play it back.

Every 5 seconds, 
If the clock is set to S/PDIF< and the clock isn't locked, clock is set to Internal.  
If the clock is set to internal and S/PDIF input has some input, it will set the clock to S/PDIF.

Don't use this code, it's likely dangerous, you might discover a UDP service and get bad input.
