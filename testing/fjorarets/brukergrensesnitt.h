#ifndef __INCLUDE_BRUKERGRENSESNITT_H__
#define __INCLUDE_BRUKERGRENSESNITT_H__

// Ser om knapper blir trigget, og setter tilhørende lys høyt.
void set_button_lights();

// Sjekker om nødstoppknappen er aktivert. Setter lampen, nullstiller ordre og stopper heis. 
int check_estop();

// Sjekker om det er ordre over.
int any_order_above();

// Sjekker om det er ordre under.
int any_order_below();

// Nullstiller alle etasjelys.
void reset_all_orders();

#endif // #ifndef __INCLUDE_BRUKERGRENSESNITT_H__
