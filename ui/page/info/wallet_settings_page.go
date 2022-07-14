package info

import (
	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletSettingsPageID = "WalletSettings"

type WalletSettingsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	wallet *dcrlibwallet.Wallet

	changePass, rescan, deleteWallet *decredmaterial.Clickable

	chevronRightIcon *decredmaterial.Icon
	backButton       decredmaterial.IconButton
	infoButton       decredmaterial.IconButton
}

func NewWalletSettingsPage(l *load.Load, wal *dcrlibwallet.Wallet) *WalletSettingsPage {
	pg := &WalletSettingsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(WalletSettingsPageID),
		wallet:           wal,
		changePass:       l.Theme.NewClickable(false),
		rescan:           l.Theme.NewClickable(false),
		deleteWallet:     l.Theme.NewClickable(false),

		chevronRightIcon: decredmaterial.NewIcon(l.Theme.Icons.ChevronRight),
	}

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) OnNavigatedTo() {
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) Layout(gtx layout.Context) layout.Dimensions {

	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrSettings),
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !pg.wallet.IsWatchingOnlyWallet() {
							return pg.changePassphrase()(gtx)
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(pg.debug()),
					layout.Rigid(pg.dangerZone()),
				)
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx, body)
	}
	return pg.layoutDesktop(gtx, body)
}

func (pg *WalletSettingsPage) layoutDesktop(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return components.UniformPadding(gtx, body)
}

func (pg *WalletSettingsPage) layoutMobile(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return components.UniformMobile(gtx, false, false, body)
}

func (pg *WalletSettingsPage) changePassphrase() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrGeneral),
			pg.bottomSectionLabel(pg.changePass, values.String(values.StrSpendingPassword)))
	}
}

func (pg *WalletSettingsPage) debug() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrDebug),
			pg.bottomSectionLabel(pg.rescan, values.String(values.StrRescanBlockchain)))
	}
}

func (pg *WalletSettingsPage) dangerZone() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrDangerZone),
			pg.bottomSectionLabel(pg.deleteWallet, values.String(values.StrRemoveWallet)))
	}
}

func (pg *WalletSettingsPage) pageSections(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	dims := func(gtx layout.Context, title string, body layout.Widget) D {
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.Body2(title)
							txt.Color = pg.Theme.Color.GrayText2
							return txt.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							if title == values.String(values.StrGeneral) {
								pg.infoButton.Inset = layout.UniformInset(values.MarginPadding0)
								pg.infoButton.Size = values.MarginPadding20
								return layout.E.Layout(gtx, pg.infoButton.Layout)
							}
							return D{}
						}),
					)
				}),
				layout.Rigid(body),
			)
		})
	}

	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			return dims(gtx, title, body)
		})
	})
}

func (pg *WalletSettingsPage) bottomSectionLabel(clickable *decredmaterial.Clickable, title string) layout.Widget {
	return func(gtx C) D {
		return clickable.Layout(gtx, func(gtx C) D {
			textLabel := pg.Theme.Body1(title)
			if title == values.String(values.StrRemoveWallet) {
				textLabel.Color = pg.Theme.Color.Danger
			}
			return layout.Inset{
				Top:    values.MarginPadding15,
				Bottom: values.MarginPadding4,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(textLabel.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							pg.chevronRightIcon.Color = pg.Theme.Color.Gray1
							return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
						})
					}),
				)
			})
		})
	}
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) HandleUserInteractions() {
	for pg.changePass.Clicked() {
		currentSpendingPasswordModal := modal.NewPasswordModal(pg.Load).
			Title(values.String(values.StrChangeSpendingPass)).
			Hint("Current spending password").
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
				go func() {
					err := pg.wallet.UnlockWallet([]byte(password))
					if err != nil {
						pm.SetError(err.Error())
						pm.SetLoading(false)
						return
					}
					pg.wallet.LockWallet()
					pm.Dismiss()

					// change password
					newSpendingPasswordModal := modal.NewCreatePasswordModal(pg.Load).
						Title(values.String(values.StrChangeSpendingPass)).
						EnableName(false).
						PasswordHint("New spending password").
						ConfirmPasswordHint("Confirm new spending password").
						PasswordCreated(func(walletName, newPassword string, m *modal.CreatePasswordModal) bool {
							go func() {
								err := pg.WL.MultiWallet.ChangePrivatePassphraseForWallet(pg.wallet.ID, []byte(password),
									[]byte(newPassword), dcrlibwallet.PassphraseTypePass)
								if err != nil {
									m.SetError(err.Error())
									m.SetLoading(false)
									return
								}
								pg.Toast.Notify("Spending password updated")
								m.Dismiss()
							}()
							return false
						})
					pg.ParentWindow().ShowModal(newSpendingPasswordModal)

				}()

				return false
			})
		pg.ParentWindow().ShowModal(currentSpendingPasswordModal)
		break
	}

	for pg.rescan.Clicked() {
		go func() {
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrRescanBlockchain)).
				Body("Rescanning may help resolve some balance errors. This will take some time, as it scans the entire"+
					" blockchain for transactions").
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(values.String(values.StrRescan), func(isChecked bool) bool {
					err := pg.WL.MultiWallet.RescanBlocks(pg.wallet.ID)
					if err != nil {
						if err.Error() == dcrlibwallet.ErrNotConnected {
							pg.Toast.NotifyError(values.String(values.StrNotConnected))
							return true
						}
						pg.Toast.NotifyError(err.Error())
						return true
					}
					msg := values.String(values.StrRescanProgressNotification)
					pg.Toast.Notify(msg)
					return true
				})

			pg.ParentWindow().ShowModal(info)
		}()
		break
	}

	for pg.deleteWallet.Clicked() {
		warningMsg := "Make sure to have the seed word backed up before removing the wallet"
		if pg.wallet.IsWatchingOnlyWallet() {
			warningMsg = "The watch-only wallet will be removed from your app"
		}
		confirmRemoveWalletModal := modal.NewInfoModal(pg.Load)
		confirmRemoveWalletModal.Title(values.String(values.StrRemoveWallet)).
			Body(warningMsg).
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButtonStyle(pg.Load.Theme.Color.Surface, pg.Load.Theme.Color.Danger).
			PositiveButton(values.String(values.StrRemove), func(isChecked bool) bool {
				walletDeleted := func() {
					if pg.WL.MultiWallet.LoadedWalletsCount() > 0 {
						pg.Toast.Notify(values.String(values.StrWalletRemoved))
						pg.ParentNavigator().CloseCurrentPage()
					} else {
						pg.ParentWindow().CloseAllPages()
					}
				}

				if pg.wallet.IsWatchingOnlyWallet() {
					confirmRemoveWalletModal.SetLoading(true)
					go func() {
						// no password is required for watching only wallets.
						err := pg.WL.MultiWallet.DeleteWallet(pg.wallet.ID, nil)
						if err != nil {
							pg.Toast.NotifyError(err.Error())
							confirmRemoveWalletModal.SetLoading(false)
						} else {
							walletDeleted()
							confirmRemoveWalletModal.Dismiss()
						}
					}()
					return false
				}

				walletPasswordModal := modal.NewPasswordModal(pg.Load).
					Title(values.String(values.StrConfirmToRemove)).
					NegativeButton(values.String(values.StrCancel), func() {}).
					PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
						go func() {
							err := pg.WL.MultiWallet.DeleteWallet(pg.wallet.ID, []byte(password))
							if err != nil {
								pm.SetError(err.Error())
								pm.SetLoading(false)
								return
							}

							walletDeleted()
							pm.Dismiss() // calls RefreshWindow.
						}()
						return false
					})
				pg.ParentWindow().ShowModal(walletPasswordModal)
				return false
			})
		pg.ParentWindow().ShowModal(confirmRemoveWalletModal)
	}

	if pg.infoButton.Button.Clicked() {
		info := modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrSpendingPassword)).
			Body(values.String(values.StrSpendingPasswordInfo)).
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
				return true
			})
		pg.ParentWindow().ShowModal(info)
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) OnNavigatedFrom() {}