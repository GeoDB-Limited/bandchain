type justify =
  | Center
  | Between
  | Right;

module Styles = {
  open Css;
  open Belt.Option;

  let row =
    style([
      display(`flex),
      flex(`num(1.)),
      width(`percent(100.)),
      margin4(
        ~top=`px(0),
        ~right=`px((-1) * Spacing.unit),
        ~left=`px((-1) * Spacing.unit),
        ~bottom=`px(0),
      ),
    ]);

  let justify =
    mapWithDefault(
      _,
      style([justifyContent(`normal)]),
      fun
      | Center => style([justifyContent(`center)])
      | Between => style([justifyContent(`spaceBetween)])
      | Right => style([justifyContent(`right)]),
    );

  let wrap = style([flexWrap(`wrap)]);
};

[@react.component]
let make = (~justify=?, ~alignItems=?, ~wrap=?, ~children) => {
  <div
    className={Cn.make([
      Styles.row,
      Styles.justify(justify),
      Styles.wrap->Cn.ifTrue(wrap->Belt.Option.getWithDefault(false)),
      // Perhaps the above props should just be a direct map like below...
      Css.style([
        alignItems->Belt.Option.getWithDefault(`center)->Css.alignItems,
      ]),
    ])}>
    children
  </div>;
};